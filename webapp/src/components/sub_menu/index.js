import {connect} from 'react-redux';

import {isSubmenuDisplayed} from 'selectors';

import SubMenu from './sub_menu';

const mapStateToProps = (state) => ({
    display: isSubmenuDisplayed(state),
});

export default connect(mapStateToProps)(SubMenu);
