import {
  combineReducers,
  configureStore,
  PreloadedState,
} from '@reduxjs/toolkit';
import { loaderReducer } from './slices/loaderSlice';
import { sessionReducer } from './slices/sessionSlice';
// import { examManagerReducer } from './slices/examManagerSlice';
// import { questionManagerReducer } from './slices/questionManagerSlice';
// import { askFormReducer } from './slices/askFormSlice';
// import { examRecordManagerReducer } from './slices/examRecordManagerSlice';

// Create the root reducer independently to obtain the RootState type
const rootReducer = combineReducers({
  loader: loaderReducer,
  session: sessionReducer,
  // examManager: examManagerReducer,
  // questionManager: questionManagerReducer,
  // askForm: askFormReducer,
  // examRecordManager: examRecordManagerReducer,
});

export type RootState = ReturnType<typeof rootReducer>;
export function setupStore(preloadedState?: PreloadedState<RootState>) {
  return configureStore({
    reducer: rootReducer,
    preloadedState,
  });
}
export type AppStore = ReturnType<typeof setupStore>;
export type AppDispatch = AppStore['dispatch'];
