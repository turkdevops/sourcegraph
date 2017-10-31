import escapeRegexp from 'escape-string-regexp'
import * as H from 'history'
import * as path from 'path'
import * as React from 'react'
import { matchPath } from 'react-router'
import 'rxjs/add/operator/catch'
import 'rxjs/add/operator/map'
import { Subscription } from 'rxjs/Subscription'
import { routes } from '../routes'
import { fetchSearchScopes } from './backend'

interface Props {
    location: H.Location

    /**
     * The query of the active search scope, or undefined if it's still loading
     */
    value?: string

    /**
     * Called when there is a change to the query provided by the active search scope.
     */
    onChange: (query: string) => void
}

interface ISearchScope {
    name: string
    value: string
}

/** The subset of State that is persisted to localStorage */
interface PersistableState {
    /** All fetched search scopes */
    remoteScopes?: ISearchScope[]
}

/** Data that is persisted to localStorage but NOT in the component state. */
interface PersistedState extends PersistableState {
    /** The value of the last-active scope. */
    lastScopeValue?: string
}

interface State extends PersistableState { }

export class SearchScope extends React.Component<Props, State> {
    private static REMOTE_SCOPES_STORAGE_KEY = 'SearchScope/remoteScopes'
    private static LAST_SCOPE_STORAGE_KEY = 'SearchScope/lastScope'

    private subscriptions = new Subscription()

    private selectElement: HTMLSelectElement | null

    constructor(props: Props) {
        super(props)

        const savedState = this.loadFromLocalStorage()
        this.state = { remoteScopes: savedState.remoteScopes }

        if (typeof savedState.lastScopeValue === 'string' && !this.props.value) {
            props.onChange(savedState.lastScopeValue)
        }

        this.subscriptions.add(
            fetchSearchScopes()
                .catch(err => {
                    console.error(err)
                    return []
                })
                .map((scopes: GQL.ISearchScope2[]) => ({ remoteScopes: scopes }))
                .subscribe(
                newState =>
                    this.setState(newState, () => {
                        this.saveToLocalStorage()

                        // Default to 1st remote scope if none.
                        if (this.props.value === undefined) {
                            this.props.onChange(newState.remoteScopes[0].value)
                        }
                    }),
                err => console.error(err)
                )
        )
    }

    public componentDidMount(): void {
        const value = this.selectElement!.value
        if (this.props.value !== undefined && this.props.value !== value) {
            this.props.onChange(value)
            this.saveToLocalStorage()
        }
    }

    public componentWillReceiveProps(newProps: Props): void {
        if (newProps.value !== undefined && newProps.value !== this.props.value) {
            this.props.onChange(newProps.value)
            this.saveToLocalStorage(newProps)
        }
    }

    public render(): JSX.Element | null {
        const scopes = this.getScopes()

        return (
            <div
                className='search-scope2'
            >
                <select
                    className='search-scope2__select ui-text-box'
                    onChange={this.onChange}
                    value={this.props.value}
                    ref={e => this.selectElement = e}
                    title='Search scope'
                >
                    {scopes.map((scope, i) =>
                        <option key={i} value={scope.value}>{scope.name}</option>
                    )}
                </select>
            </div>
        )
    }

    private onChange: React.ChangeEventHandler<HTMLSelectElement> = event => {
        this.props.onChange(event.currentTarget.value)
    }

    private getScopes(): ISearchScope[] {
        const allScopes: ISearchScope[] = []

        if (this.state.remoteScopes) {
            allScopes.push(... this.state.remoteScopes)
        }

        allScopes.push(...this.getScopesForCurrentRoute())

        // If the active scope isn't in the list, then add a new custom entry for it.
        if (this.props.value !== undefined && !allScopes.some(({ value }) => value === this.props.value)) {
            allScopes.push({
                name: 'Custom',
                value: this.props.value,
            })
        }

        return allScopes
    }

    /**
     * Returns contextual scopes for the current route (such as "This Repository" and
     * "This Directory").
     */
    private getScopesForCurrentRoute(): ISearchScope[] {
        const scopes: ISearchScope[] = []

        // This is basically a programmatical <Switch> with <Route>s
        // see https://reacttraining.com/react-router/web/api/matchPath
        // and https://reacttraining.com/react-router/web/example/sidebar
        for (const route of routes) {
            const match = matchPath<{ repoRev?: string, filePath?: string }>(this.props.location.pathname, route)
            if (match) {
                switch (match.path) {
                    case '/:repoRev+': {
                        // Repo page
                        const [repoPath] = match.params.repoRev!.split('@')
                        scopes.push({
                            name: `This repository (${path.basename(repoPath)})`,
                            value: `repo:^${escapeRegexp(repoPath)}$`,
                        })
                        break
                    }
                    case '/:repoRev+/-/blob/:filePath+': {
                        // Blob/tree page
                        const [repoPath] = match.params.repoRev!.split('@')

                        scopes.push({
                            name: `This repository (${path.basename(repoPath)})`,
                            value: `repo:^${escapeRegexp(repoPath)}$`,
                        })

                        if (match.params.filePath) {
                            const dirname = path.dirname(match.params.filePath)
                            if (dirname !== '.') {
                                scopes.push({
                                    name: `This directory (${path.basename(dirname)})`,
                                    value: `repo:^${escapeRegexp(repoPath)}$ file:^${escapeRegexp(dirname)}/`,
                                })
                            }
                        }
                        break
                    }
                }
                break
            }
        }

        return scopes
    }

    private saveToLocalStorage(props: Props = this.props): void {
        const writeItem = (key: string, data: any): void => {
            if (data !== undefined && data !== null) {
                localStorage.setItem(key, JSON.stringify(data))
            } else {
                localStorage.removeItem(key)
            }
        }

        writeItem(SearchScope.REMOTE_SCOPES_STORAGE_KEY, this.state.remoteScopes)
        writeItem(SearchScope.LAST_SCOPE_STORAGE_KEY, props.value)
    }

    private loadFromLocalStorage(): PersistedState {
        const readItem = <T extends {}>(key: string, validate: (data: T) => boolean): T | undefined => {
            const raw = localStorage.getItem(key)
            if (raw === null) { return undefined }

            try {
                const data = JSON.parse(raw)
                if (data !== undefined && data !== null && validate(data)) { return data }
            } catch (err) { /* noop */ }

            // Else invalid data.
            localStorage.removeItem(key)
            return undefined
        }

        const validate = (data: ISearchScope): boolean => typeof data.name === 'string' && typeof data.value === 'string'
        const remoteScopes = readItem<ISearchScope[]>(SearchScope.REMOTE_SCOPES_STORAGE_KEY, data => data.every(validate))
        const lastScopeValue = readItem<string>(SearchScope.LAST_SCOPE_STORAGE_KEY, s => typeof s === 'string')
        return { remoteScopes, lastScopeValue }
    }
}
