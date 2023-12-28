import { Provider } from 'react-redux';
import { ChakraProvider } from '@chakra-ui/react';
import axios from 'axios';
import { HotkeysProvider } from 'react-hotkeys-hook';
import { setupStore } from './store/store';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import theme from './themes/theme';
import Root from './routes/Root';
import Home from './routes/Home';
import SignUp from './routes/SignUp';
import SignIn from './routes/SignIn';
import Word from './routes/Word';

axios.defaults.baseURL = `${import.meta.env.VITE_BACKEND_URL}/api`;
axios.defaults.withCredentials = true;

// JWT
axios.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');

    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

function App() {
  const store = setupStore({});
  const router = createBrowserRouter([
    {
      path: '/',
      element: <Root />,
      children: [
        {
          path: '/',
          element: <Home />,
        },
        {
          path: '/signup',
          element: <SignUp />,
        },
        {
          path: '/signin',
          element: <SignIn />,
        },
        {
          path: '/restricted/word',
          element: <Word />,
        },
      ],
    },
  ]);

  return (
    <ChakraProvider theme={theme}>
      <Provider store={store}>
        {/*
            HotkeysProvider 設定快捷鍵的 scope 
            'card': 給"單字卡"頁面使用
        */}
        <HotkeysProvider initiallyActiveScopes={['card']}>
          <RouterProvider router={router} />
        </HotkeysProvider>
      </Provider>
    </ChakraProvider>
  );
}

export default App;
