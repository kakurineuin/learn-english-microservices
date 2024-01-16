import { Provider } from 'react-redux';
import { ChakraProvider } from '@chakra-ui/react';
import axios from 'axios';
import { HotkeysProvider } from 'react-hotkeys-hook';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { setupStore } from './store/store';
import theme from './themes/theme';
import Root from './routes/Root';
import Home from './routes/Home';
import SignUp from './routes/SignUp';
import SignIn from './routes/SignIn';
import WordManager from './routes/WordManager';
import FavoriteWordMeaningManager from './routes/FavoriteWordMeaningManager';
import WordCard from './routes/WordCard';
import ExamManager from './routes/ExamManager';
import QuestionManager from './routes/QuestionManager';
import StartExam from './routes/StartExam';
import ExamRecordOverview from './routes/ExamRecordOverview';

axios.defaults.baseURL = `${import.meta.env.VITE_BACKEND_URL}/api`;
axios.defaults.withCredentials = true;
axios.defaults.xsrfCookieName = 'XSRF-TOKEN';
axios.defaults.xsrfHeaderName = 'X-XSRF-TOKEN';

// JWT
axios.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');

    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => Promise.reject(error),
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
          path: '/restricted/exam',
          element: <ExamManager />,
        },
        {
          path: '/restricted/exam/:examId',
          element: <QuestionManager />,
        },
        {
          path: '/restricted/exam/:examId/start',
          element: <StartExam />,
        },
        {
          path: '/restricted/exam/:examId/record/overview',
          element: <ExamRecordOverview />,
        },
        {
          path: '/restricted/word',
          element: <WordManager />,
        },
        {
          path: '/restricted/word/favorite',
          element: <FavoriteWordMeaningManager />,
        },
        {
          path: '/restricted/word/card',
          element: <WordCard />,
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
