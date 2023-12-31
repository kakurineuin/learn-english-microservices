import { useMemo, useCallback } from 'react';
import { VStack, Button, useToast, Center } from '@chakra-ui/react';
import { CellProps, Column } from 'react-table';
import { DateTime } from 'luxon';
import { useNavigate } from 'react-router-dom';
import { useAppSelector, useAppDispatch } from '../../store/hooks';
import { Exam } from '../../models/Exam';
import { queryExams, deleteExam } from '../../store/slices/examManagerSlice';
import PaginationTable from '../PaginationTable';
import UpdateExamDialog from './UpdateExamDialog';
import ConfirmDialog from '../ConfirmDialog';

function ExamList() {
  const user = useAppSelector((state) => state.session.user);
  const userRole = user!.role;
  const saveExamTimes = useAppSelector(
    (state) => state.examManager.saveExamTimes,
  );
  const refreshId = saveExamTimes;
  const dispatch = useAppDispatch();
  const toast = useToast();
  const navigate = useNavigate();
  const columns = useMemo(
    () =>
      [
        {
          Header: 'Topic',
          accessor: 'topic',
        },
        {
          Header: 'Description',
          accessor: 'description',
          Cell: ({ value }: CellProps<Exam, string>) => <pre>{value}</pre>,
        },
        {
          Header: 'Private / Public',
          accessor: 'isPublic',
          Cell: ({ value }: CellProps<Exam, boolean>) => (
            <Center>{value ? 'Public' : 'Private'}</Center>
          ),
        },
        {
          Header: 'Updated At',
          accessor: 'updatedAt',
          Cell: ({ value }: CellProps<Exam, string>) => (
            <Center>
              {DateTime.fromISO(value).toFormat('yyyy/LL/dd HH:mm:ss')}
            </Center>
          ),
        },
        {
          Header: 'Edit',
          accessor: '_id',
          width: 50,
          Cell: ({ row }: CellProps<Exam, string>) => (
            <VStack spacing={2}>
              <UpdateExamDialog
                key={row.values._id}
                data={row.values as Exam}
              />
              <Button
                colorScheme="teal"
                variant="outline"
                size="sm"
                onClick={() => {
                  navigate(`/restricted/exam/${row.values._id}`);
                }}
              >
                題目
              </Button>
              <ConfirmDialog
                button={{
                  colorScheme: 'red',
                  variant: 'outline',
                  size: 'sm',
                  text: '刪除',
                }}
                header="刪除確認"
                yesCallback={() => {
                  dispatch(deleteExam(row.values._id))
                    .unwrap()
                    .then(() => {
                      toast({
                        title: '刪除成功。',
                        status: 'success',
                        isClosable: true,
                        position: 'top',
                        variant: 'subtle',
                      });
                    })
                    .catch((message) => {
                      toast({
                        title: '刪除失敗。',
                        description: message,
                        status: 'error',
                        isClosable: true,
                        position: 'top',
                        variant: 'subtle',
                      });
                    });
                }}
              >
                確定要刪除此份測驗 &quot;{row.values.topic}&quot;？
              </ConfirmDialog>
            </VStack>
          ),
        },
      ] as Column<Exam>[],
    [dispatch, navigate, toast],
  );
  const total = useAppSelector((state) => state.examManager.total);
  const pageCount = useAppSelector((state) => state.examManager.pageCount);
  const data = useAppSelector((state) => state.examManager.exams);
  const fetchData = useCallback(
    (pageIndex: number, pageSize: number) => {
      dispatch(queryExams({ pageIndex, pageSize }))
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

  // 管理者顯示全部欄位，一般使用者不用顯示 isPublic 欄位
  const hiddenColumns = useMemo(
    () => (userRole === 'admin' ? [] : ['isPublic']),
    [userRole],
  );

  return (
    <PaginationTable
      columns={columns}
      hiddenColumns={hiddenColumns}
      data={data}
      fetchData={fetchData}
      pageCount={pageCount}
      total={total}
      refreshId={refreshId}
    />
  );
}

export default ExamList;
