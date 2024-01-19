import { useMemo, useCallback } from 'react';
import { useToast, Center } from '@chakra-ui/react';
import { CellProps, Column } from 'react-table';
import { DateTime } from 'luxon';
import { useAppSelector, useAppDispatch } from '../../store/hooks';
import { UserHistory } from '../../models/UserHistory';
import { queryUserHistories } from '../../store/slices/userHistoryManagerSlice';
import PaginationTable from '../PaginationTable';

function UserHistoryList() {
  const dispatch = useAppDispatch();
  const toast = useToast();
  const columns = useMemo(
    () =>
      [
        {
          Header: 'Created At',
          accessor: 'createdAt',
          Cell: ({ value }: CellProps<UserHistory, string>) => (
            <Center>
              {DateTime.fromISO(value).toFormat('yyyy/LL/dd HH:mm:ss')}
            </Center>
          ),
        },
        {
          Header: 'Username',
          accessor: 'username',
          Cell: ({ value }: CellProps<UserHistory, string>) =>
            value || '未登入',
        },
        {
          Header: 'Role',
          accessor: 'role',
        },
        {
          Header: 'Method',
          accessor: 'method',
        },
        {
          Header: 'Path',
          accessor: 'path',
        },
      ] as Column<UserHistory>[],
    [],
  );
  const total = useAppSelector((state) => state.userHistoryManager.total);
  const pageCount = useAppSelector(
    (state) => state.userHistoryManager.pageCount,
  );
  const data = useAppSelector(
    (state) => state.userHistoryManager.userHistories,
  );
  const fetchData = useCallback(
    (pageIndex: number, pageSize: number) => {
      dispatch(queryUserHistories({ pageIndex, pageSize }))
        .unwrap()
        .catch((message) => {
          toast({
            title: '查詢失敗。',
            description: message,
            status: 'error',
            isClosable: true,
            position: 'top',
            variant: 'subtle',
          });
        });
    },
    [dispatch, toast],
  );

  return (
    <PaginationTable
      columns={columns}
      data={data}
      fetchData={fetchData}
      pageCount={pageCount}
      total={total}
      refreshId={1}
    />
  );
}

export default UserHistoryList;
