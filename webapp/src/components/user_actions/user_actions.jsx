import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';

export default class UserActionsComponent extends React.PureComponent {
    static propTypes = {
        openRootModal: PropTypes.func.isRequired,
        theme: PropTypes.object.isRequired,
    }

    onClick = () => {
        this.props.openRootModal();
    }

    render() {
        const style = getStyle(this.props.theme);

        return (
            <div>
                <FormattedMessage
                    id='useractions.demo'
                    defaultMessage='Demo Plugin: '
                />
                <button
                    style={style.button}
                    onClick={this.onClick}
                >
                    <FormattedMessage
                        id='useractions.action'
                        defaultMessage='Action'
                    />
                </button>
            </div>
        );
    }
}

const getStyle = (theme) => ({
    button: {
        color: theme.buttonColor,
        backgroundColor: theme.buttonBg,
    },
});
