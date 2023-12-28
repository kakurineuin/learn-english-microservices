import { Outlet } from 'react-router-dom';
import Loader from '../components/Loader';
import { useAppSelector } from '../store/hooks';
import Menus from '../components/Menus';

function Root() {
  const isLoading = useAppSelector((state) => state.loader.isLoading);

  return (
    <>
      {isLoading && <Loader />}
      <div
        style={{
          minHeight: '100vh',
          display: 'flex',
          flexDirection: 'column',
          backgroundImage: 'url("/assets/bg_1920x1080.png")',
        }}
      >
        <Menus />

        <div
          style={{
            flex: 1,
          }}
          className="h-full m-2 backdrop-blur-md bg-gray-950/60 text-white border-solid border-2 border-gray-600 rounded-lg p-4"
        >
          <Outlet />
        </div>
      </div>
    </>
  );
}

export default Root;
