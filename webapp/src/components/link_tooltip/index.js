import {connect} from 'react-redux';

import {isEnabled} from 'selectors';

import LinkTooltip from './link_tooltip';

const mapStateToProps = (state) => ({
    enabled: isEnabled(state),
});

export default connect(mapStateToProps)(LinkTooltip);
