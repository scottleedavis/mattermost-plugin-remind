import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';

// LeftSidebarHeader is a pure component, later connected to the Redux store so as to
// show the plugin's enabled / disabled status.
export default class LeftSidebarHeader extends React.PureComponent {
    static propTypes = {
        enabled: PropTypes.bool.isRequired,
    }

    render() {
        const iconStyle = {
            display: 'inline-block',
            margin: '0 7px 0 1px',
        };
        const style = {
            margin: '.5em 0 .5em',
            padding: '0 12px 0 15px',
            color: 'rgba(255,255,255,0.6)',
        };

        return (
            <div style={style}>
                <i
                    className='icon fa fa-plug'
                    style={iconStyle}
                />
                <FormattedMessage
                    id='sidebar.demo'
                    defaultMessage='Demo Plugin:'
                />
                {' '}
                {this.props.enabled ?
                    <span>
                        <FormattedMessage
                            id='sidebar.enabled'
                            defaultMessage='Enabled'
                        />
                    </span> :
                    <span>
                        <FormattedMessage
                            id='sidebar.disabled'
                            defaultMessage='Disabled'
                        />
                    </span>
                }
            </div>
        );
    }
}
