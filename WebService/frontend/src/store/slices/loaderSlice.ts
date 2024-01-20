import { createSlice } from '@reduxjs/toolkit';

type LoaderState = {
  isLoading: boolean;
};

const initialState: LoaderState = {
  isLoading: false,
};

export const loaderSlice = createSlice({
  name: 'loader',
  initialState,
  reducers: {
    toggleLoading: (state) => {
      state.isLoading = !state.isLoading;
    },
  },
});

export const loaderActions = loaderSlice.actions;
export const loaderReducer = loaderSlice.reducer;
