import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {openRootModal} from 'actions';

import UserActions from './user_actions';

const mapDispatchToProps = (dispatch) => bindActionCreators({
    openRootModal,
}, dispatch);

export default connect(null, mapDispatchToProps)(UserActions);
