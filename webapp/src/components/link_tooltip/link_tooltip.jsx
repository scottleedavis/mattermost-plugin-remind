import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';

export default class LinkTooltip extends React.PureComponent {
    static propTypes = {
        href: PropTypes.string.isRequired,
    }

    render() {
        if (!this.props.href.includes('www.test.com')) {
            return null;
        }

        return (
            <div
                style={parentDivStyles}
            >
                <i
                    style={iconStyles}
                    className='icon fa fa-plug'
                />
                <FormattedMessage
                    id='tooltip.message'
                    defaultMessage='This is a custom tooltip from the Demo Plugin'
                />
            </div>
        );
    }
}

const parentDivStyles = {
    backgroundColor: '#ffffff',
    borderRadius: '4px',
    boxShadow: 'rgba(61, 60, 64, 0.1) 0px 17px 50px 0px, rgba(61, 60, 64, 0.1) 0px 12px 15px 0px',
    fontSize: '14px',
    marginTop: '10px',
    padding: '10px 15px 15px',
};

const iconStyles = {
    paddingRight: '5px',
};
