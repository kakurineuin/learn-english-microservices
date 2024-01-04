import { useState, useCallback, useEffect } from 'react';
import { Button, Box, Container, Divider, useToast } from '@chakra-ui/react';
import { motion } from 'framer-motion';
import { useNavigate, useParams } from 'react-router-dom';
import axios, { AxiosError } from 'axios';
import { fadeIn } from '../utils/motion';
import { Exam } from '../models/Exam';
import { Question } from '../models/Question';
import AskForm from '../components/exam/AskForm';
import PageHeading from '../components/PageHeading';
import ShowText from '../components/ShowText';
import { useAppDispatch } from '../store/hooks';
import { loaderActions } from '../store/slices/loaderSlice';

const refreshPage = () => {
  // 判斷是否已在瀏覽器環境
  if (typeof window !== 'undefined') {
    window.location.reload();
  }
};

function StartExam() {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const toast = useToast();
  const { examId } = useParams();
  const [exam, setExam] = useState<Exam | null>(null);
  const [questions, setQuestions] = useState<Question[]>([]);
  const [isDisabledGoExamRecord, setIsDisabledGoExamRecord] = useState(true);
  const [isShowDescription, setIsShowDescription] = useState(false);

  const onCreateExamRecordHandler = useCallback(() => {
    setIsDisabledGoExamRecord(false);
  }, []);

  const goExamRecordClickHandler = () => {
    navigate(`/restricted/exam/${examId}/record`);
  };

  useEffect(() => {
    const queryRandomQuestions = async () => {
      dispatch(loaderActions.toggleLoading());

      try {
        const response = await axios.get(`/restricted/exam/${examId}/start`);
        setExam(response.data.exam);
        setQuestions(response.data.questions);
      } catch (err) {
        const errorMessage = axios.isAxiosError(err)
          ? (err as AxiosError<{ message: string }, any>).response!.data.message
          : '系統發生錯誤！';
        toast({
          title: '查詢失敗。',
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

    queryRandomQuestions();
  }, [dispatch, toast, examId]);

  return (
    <Container maxW="container.xl" mt="3">
      {exam && (
        <motion.div
          variants={fadeIn('down', 'tween', 0.2, 1)}
          initial="hidden"
          animate="show"
          onAnimationComplete={() => setIsShowDescription(true)}
        >
          <PageHeading title={exam.topic} />
          {isShowDescription && <ShowText>{exam.description}</ShowText>}
          <Box textAlign="right">
            <Button colorScheme="teal" variant="outline" onClick={refreshPage}>
              再測驗一次
            </Button>

            <Button
              colorScheme="blue"
              variant="outline"
              ml="3"
              isDisabled={isDisabledGoExamRecord}
              onClick={goExamRecordClickHandler}
            >
              成績紀錄
            </Button>
          </Box>
        </motion.div>
      )}
      <Divider my={5} />
      {questions.length > 0 && (
        <AskForm
          questions={questions}
          onCreateExamRecord={onCreateExamRecordHandler}
        />
      )}
    </Container>
  );
}

export default StartExam;
