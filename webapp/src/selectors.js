import {id as pluginId} from './manifest';

const getPluginState = (state) => state['plugins-' + pluginId] || {};

export const isEnabled = (state) => getPluginState(state).enabled;

export const isSubmenuDisplayed = (state) => getPluginState(state).display;

export const isRootModalVisible = (state) => getPluginState(state).rootModalVisible;
