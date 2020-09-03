import * as React from 'react'
import classnames from 'classnames'

interface Props {
    /** The initial value. */
    value?: boolean

    /** The DOM ID of the element. */
    id?: string

    /**
     * Called when the user changes the input's value.
     */
    onToggle?: (value: boolean) => void

    /**
     * Called when the user hovers over the toggle.
     */
    onHover?: (value: boolean) => void

    /** The title attribute (tooltip). */
    title?: string

    disabled?: boolean
    tabIndex?: number
    className?: string
}

/** A toggle switch input component. */
export const Toggle: React.FunctionComponent<Props> = ({
    disabled,
    className,
    id,
    title,
    value,
    tabIndex,
    onToggle,
    onHover,
}) => {
    function onClick(): void {
        if (!disabled && onToggle) {
            onToggle(!value)
        }
    }

    function onMouseOver(): void {
        if (onHover) {
            onHover(!value)
        }
    }

    return (
        <button
            type="button"
            className={classnames('toggle', className, {})}
            id={id}
            title={title}
            value={value ? 1 : 0}
            onClick={onClick}
            tabIndex={tabIndex}
            onMouseOver={onMouseOver}
            disabled={disabled}
        >
            <span
                className={classnames('toggle__bar', {
                    'toggle__bar--on': value,
                    'toggle__bar--disabled': disabled,
                })}
            />
            <span
                className={classnames('toggle__bar-shadow', {
                    'toggle__bar-shadow--on': value,
                })}
            />
            <span
                className={classnames('toggle__knob', {
                    'toggle__knob--on': value,
                    'toggle__knob--disabled': disabled,
                })}
            />
        </button>
    )
}
