import {connect} from 'react-redux';

import {isEnabled} from 'selectors';

import SubMenu from './sub_menu';

const mapStateToProps = (state) => ({
    enabled: isEnabled(state),
});

export default connect(mapStateToProps)(SubMenu);
