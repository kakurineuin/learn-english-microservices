import { useEffect, useMemo, useState } from 'react';
import { VStack, Button, Text, useToast, Center } from '@chakra-ui/react';
import axios, { AxiosError } from 'axios';
import { motion } from 'framer-motion';
import { slideIn, fadeIn } from '../../utils/motion';
import AskQuestion from './AskQuestion';
import { loaderActions } from '../../store/slices/loaderSlice';
import { useAppSelector, useAppDispatch } from '../../store/hooks';
import { Question } from '../../models/Question';

type Props = {
  questions: Question[];
  onCreateExamRecord: () => void;
};

function AskForm({ questions, onCreateExamRecord }: Props) {
  const maxScore = questions.length;
  const [isValidate, setIsValidate] = useState(false);
  const dispatch = useAppDispatch();
  const toast = useToast();
  const questionResults = useAppSelector(
    (state) => state.askForm.questionResults,
  );

  const wrongQuestionIds: string[] = useMemo(() => [], []);
  let score = 0;
  let isScored = false;

  if (questionResults.length === questions.length) {
    const ids = questionResults
      .filter(({ isSuccess }) => !isSuccess)
      .map(({ questionId }) => questionId);
    wrongQuestionIds.length = 0;
    wrongQuestionIds.push(...ids);
    score = questions.length - wrongQuestionIds.length;
    isScored = true;
  }

  useEffect(() => {
    if (!isScored) {
      return;
    }

    // 保存測驗紀錄到後端
    const createExamRecord = async () => {
      dispatch(loaderActions.toggleLoading());

      try {
        await axios.post(`/exam/${questions[0].examId}/record`, {
          score,
          wrongQuestionIds,
        });
        toast({
          title: '新增紀錄成功',
          status: 'success',
          isClosable: true,
          position: 'top',
          variant: 'subtle',
        });

        onCreateExamRecord();
      } catch (err) {
        const errorMessage = axios.isAxiosError(err)
          ? (err as AxiosError<{ message: string }, any>).response!.data.message
          : '系統發生錯誤！';
        toast({
          title: '新增紀錄失敗',
          description: errorMessage,
          status: 'error',
          isClosable: true,
          position: 'top',
          variant: 'subtle',
        });
      } finally {
        dispatch(loaderActions.toggleLoading());
      }
    };
    createExamRecord();
    window.scrollTo(0, 0);
  }, [
    dispatch,
    score,
    isScored,
    questionResults.length,
    questions,
    toast,
    wrongQuestionIds,
    onCreateExamRecord,
  ]);

  const submitHandler = () => {
    setIsValidate(true);
  };

  const childComponents = questions.map((question, index) => (
    <AskQuestion
      key={question._id}
      questionNumber={index + 1}
      question={question}
      isValidate={isValidate}
      motionVariants={slideIn('left', 'tween', 0.1 * (index + 1), 1)}
    />
  ));

  return (
    <>
      {isScored && (
        <motion.div
          variants={fadeIn('down', 'tween', 0.2, 1)}
          initial="hidden"
          animate="show"
        >
          <Center>
            <Text fontSize="100px" data-testid="score" color="blue.200">
              {score}/{maxScore}
            </Text>
          </Center>
        </motion.div>
      )}
      <VStack spacing={3} align="start">
        {childComponents}

        <Button
          colorScheme="blue"
          variant="outline"
          isDisabled={isValidate}
          onClick={submitHandler}
        >
          提交
        </Button>
      </VStack>
    </>
  );
}

export default AskForm;
