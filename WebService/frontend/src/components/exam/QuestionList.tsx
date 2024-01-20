import { useMemo, useCallback } from 'react';
import { VStack, useToast, Box, Center } from '@chakra-ui/react';
import { CellProps, Column } from 'react-table';
import { DateTime } from 'luxon';
import { v4 as uuidv4 } from 'uuid';
import { useAppSelector, useAppDispatch } from '../../store/hooks';
import { Question } from '../../models/Question';
import {
  queryQuestions,
  deleteQuestion,
} from '../../store/slices/questionManagerSlice';
import PaginationTable from '../PaginationTable';
import UpdateQuestionDialog from './UpdateQuestionDialog';
import ConfirmDialog from '../ConfirmDialog';

type Props = {
  examId: string;
};

function QuestionList({ examId }: Props) {
  const saveQuestionTimes = useAppSelector(
    (state) => state.questionManager.saveQuestionTimes,
  );
  const refreshId = saveQuestionTimes;

  const dispatch = useAppDispatch();
  const toast = useToast();
  const columns = useMemo(
    () =>
      [
        {
          Header: 'Ask',
          accessor: 'ask',
          Cell: ({ value }: CellProps<Question, string>) => <pre>{value}</pre>,
        },
        {
          Header: 'Answers',
          accessor: 'answers',
          Cell: ({ value }: CellProps<Question, string[]>) => {
            const answers = value;
            return answers.map((answer, i) => (
              <Box key={uuidv4()}>{`Answer${i + 1}: ${answer}`}</Box>
            ));
          },
        },
        {
          Header: 'Updated At',
          accessor: 'updatedAt',
          Cell: ({ value }: CellProps<Question, string>) => (
            <Center>
              {DateTime.fromISO(value).toFormat('yyyy/LL/dd HH:mm:ss')}
            </Center>
          ),
        },
        {
          Header: 'ExamId',
          accessor: 'examId',
        },
        {
          Header: 'Edit',
          accessor: '_id',
          width: 50,
          Cell: ({ row }: CellProps<Question, string>) => (
            <VStack spacing={2}>
              <UpdateQuestionDialog
                key={row.values._id}
                data={row.values as Question}
              />
              <ConfirmDialog
                button={{
                  colorScheme: 'red',
                  variant: 'outline',
                  size: 'sm',
                  text: '刪除',
                }}
                header="刪除確認"
                yesCallback={() => {
                  dispatch(deleteQuestion(row.values as Question))
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
                確定要刪除試題？
                <pre>{row.values.ask}</pre>
              </ConfirmDialog>
            </VStack>
          ),
        },
      ] as Column<Question>[],
    [dispatch, toast],
  );
  const hiddenColumns = useMemo(() => ['examId'], []);
  const total = useAppSelector((state) => state.questionManager.total);
  const pageCount = useAppSelector((state) => state.questionManager.pageCount);
  const data = useAppSelector((state) => state.questionManager.questions);
  const fetchData = useCallback(
    (pageIndex: number, pageSize: number) => {
      dispatch(queryQuestions({ examId, pageIndex, pageSize }));
    },
    [examId, dispatch],
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

export default QuestionList;
