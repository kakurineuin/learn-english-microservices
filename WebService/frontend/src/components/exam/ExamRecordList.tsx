import { useMemo, useCallback } from 'react';
import { useToast, Center } from '@chakra-ui/react';
import { CellProps, Column } from 'react-table';
import { DateTime } from 'luxon';
import { useAppSelector, useAppDispatch } from '../../store/hooks';
import { ExamRecord } from '../../models/ExamRecord';
import { queryExamRecords } from '../../store/slices/examRecordManagerSlice';
import PaginationTable from '../PaginationTable';

type Props = {
  examId: string;
};

function ExamRecordList({ examId }: Props) {
  const dispatch = useAppDispatch();
  const toast = useToast();
  const columns = useMemo(
    () =>
      [
        {
          Header: 'Created At',
          accessor: 'createdAt',
          Cell: ({ value }: CellProps<ExamRecord, string>) => (
            <Center>
              {DateTime.fromISO(value).toFormat('yyyy/LL/dd HH:mm:ss')}
            </Center>
          ),
        },
        {
          Header: 'Score',
          accessor: 'score',
          Cell: ({ value }: CellProps<ExamRecord, number>) => (
            <Center>{value}</Center>
          ),
        },
      ] as Column<ExamRecord>[],
    [],
  );
  const total = useAppSelector((state) => state.examRecordManager.total);
  const pageCount = useAppSelector(
    (state) => state.examRecordManager.pageCount,
  );
  const data = useAppSelector((state) => state.examRecordManager.examRecords);
  const fetchData = useCallback(
    (pageIndex: number, pageSize: number) => {
      dispatch(queryExamRecords({ pageIndex, pageSize, examId }))
        .unwrap()
        .catch((message) => {
          toast({
            title: '查詢失敗',
            description: message,
            status: 'error',
            isClosable: true,
            position: 'top',
            variant: 'subtle',
          });
        });
    },
    [dispatch, toast, examId],
  );

  return (
    <PaginationTable
      columns={columns}
      data={data}
      fetchData={fetchData}
      pageCount={pageCount}
      total={total}
      refreshId={0}
    />
  );
}

export default ExamRecordList;
