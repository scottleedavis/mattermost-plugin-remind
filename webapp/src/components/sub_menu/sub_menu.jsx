import React from 'react';

// import ReactDOM from 'react-dom';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';

export default class SubMenu extends React.PureComponent {
    static propTypes = {
        display: PropTypes.bool,
        post: PropTypes.object,
        postId: PropTypes.string,
    };

    constructor(props) {
        super(props);

        console.log(this.props.post); //eslint-disable-line
        console.log(this.props.postId); //eslint-disable-line
        this.state = {
            display: this.props.display || false,
        };

        // this.handleMouseEnter = this.handleMouseEnter.bind(this);
        // this.handleMouseLeave = this.handleMouseLeave.bind(this);
    }

    handleMouseEnter() {
        // this.state.display = true;
        this.setState({display: true});
        console.log("mouse enter "+this.state.display); //eslint-disable-line
        // const containerId = ReactDOM.findDOMNode(this).parentNode.parentNode.parentNode.parentNode.getAttribute('id');
        // console.log(containerId);  //eslint-disable-line
    }

    handleMouseLeave() {
        this.setState({display: false});

        // this.state.display = false;
        console.log("mouse leave "+this.state.display);  //eslint-disable-line
        // const containerId = ReactDOM.findDOMNode(this).parentNode.parentNode.parentNode.parentNode.getAttribute('id');
        // console.log(containerId);  //eslint-disable-line
    }

    render() {
        const iconStyle = {
            display: 'inline-block',
            margin: '0 7px 0 1px',
        };
        const style = {

            // margin: '.5em 0 .5em',
            // padding: '0 12px 0 15px',
            // color: 'rgba(0,0,0,0.6)',
        };
        const submenuStyle = {
            position: 'absolute',
            border: 'solid',
            height: '200px',
            width: '200px',
            right: '200px',
        };

        // console.log(this.props.post); //eslint-disable-line
        console.log(this.props.postId); //eslint-disable-line

        return (
            <div
                style={style}

                // onMouseEnter={this.handleMouseEnter}
                // onMouseLeave={this.handleMouseLeave}
            >
                {this.state.display ?
                    <ul
                        className=''
                        style={submenuStyle}
                    >
                        <li> {'display'} </li>
                    </ul> :
                    ''
                }
                <i
                    className='icon fa fa-chevron-left'
                    style={iconStyle}
                />
                <FormattedMessage
                    id='submenu.message'
                    defaultMessage='Remind me about this'
                />
            </div>
        );
    }
}
