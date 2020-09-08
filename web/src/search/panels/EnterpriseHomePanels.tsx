import * as React from 'react'
import { AuthenticatedUser } from '../../auth'
import { RecentFilesPanel } from './RecentFilesPanel'
import { RecentSearchesPanel } from './RecentSearchesPanel'
import { RepositoriesPanel } from './RepositoriesPanel'
import { SavedSearchesPanel } from './SavedSearchesPanel'
import { Observable } from 'rxjs'
import { EventLogResult } from '../backend'
import { PatternTypeProps } from '..'

interface Props extends Pick<PatternTypeProps, 'patternType'> {
    authenticatedUser: AuthenticatedUser | null
    fetchRecentSearches: (userId: string, first: number) => Observable<EventLogResult>
}

export const EnterpriseHomePanels: React.FunctionComponent<Props> = (props: Props) => (
    <div className="enterprise-home-panels container">
        <div className="row">
            <RepositoriesPanel className="enterprise-home-panels__panel col-lg-4" />
            <RecentSearchesPanel {...props} className="enterprise-home-panels__panel col-lg-8" />
        </div>
        <div className="row">
            <RecentFilesPanel className="enterprise-home-panels__panel col-lg-7" />
            <SavedSearchesPanel className="enterprise-home-panels__panel col-lg-5" />
        </div>
    </div>
)
