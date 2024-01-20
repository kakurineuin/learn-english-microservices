import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  Box,
  Container,
  Heading,
  Divider,
  Tooltip,
  useColorMode,
  useToast,
} from '@chakra-ui/react';
import { InfoIcon } from '@chakra-ui/icons';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip as ReChartTooltip,
  Legend,
  Label,
} from 'recharts';
import _ from 'lodash';
import { DateTime } from 'luxon';
import axios, { AxiosError } from 'axios';
import ExamRecordList from '../components/exam/ExamRecordList';
import WrongQuestionList from '../components/exam/WrongQuestionList';
import { Question } from '../models/Question';
import { Exam } from '../models/Exam';
import PageHeading from '../components/PageHeading';
import { useAppDispatch } from '../store/hooks';
import { loaderActions } from '../store/slices/loaderSlice';
import { AnswerWrong } from '../models/AnswerWrong';
import { ExamRecord } from '../models/ExamRecord';

interface WrongQuestion extends Question {
  answerWrongTimes: number;
}

type ExamScore = {
  _id: string;
  score: number;
};

interface ResponseData {
  startDate: string;
  exam: Exam;
  questions: Question[];
  answerWrongs: AnswerWrong[];
  examRecords: ExamRecord[];
}

function ExamRecordOverview() {
  const [startDate, setStartDate] = useState<string>('');
  const [exam, setExam] = useState<Exam | null>(null);
  const [wrongQuestions, setWrongQuestions] = useState<WrongQuestion[]>([]);
  const [examScores, setExamScores] = useState<ExamScore[]>([]);

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { colorMode, toggleColorMode } = useColorMode();
  const chartFontColor = colorMode === 'light' ? 'black' : 'white';
  const dispatch = useAppDispatch();
  const toast = useToast();
  const { examId } = useParams();

  useEffect(() => {
    const queryExamRecordOverview = async () => {
      dispatch(loaderActions.toggleLoading());

      try {
        const response = await axios.get(
          `/restricted/exam/${examId}/record/overview`,
        );
        const {
          startDate: myStartDate,
          exam: myExam,
          questions: myQuestions,
          answerWrongs: myAnswerWrongs,
          examRecords: myExamRecords,
        }: ResponseData = response.data;

        setStartDate(myStartDate);
        setExam(myExam);

        const myWrongQuestions = myQuestions.map((question) => {
          const answerWrong = myAnswerWrongs.find(
            (obj: AnswerWrong) => obj.questionId === question._id,
          );
          return {
            ...question,
            answers: [...question.answers],
            answerWrongTimes: answerWrong!.times,
          };
        });
        setWrongQuestions(myWrongQuestions);

        const groupExamRecords = _.groupBy(myExamRecords, (examRecord) =>
          DateTime.fromISO(examRecord.createdAt).toFormat('yyyy/LL/dd'),
        );

        const myExamScores: ExamScore[] = [];
        Object.entries(groupExamRecords).forEach(
          ([strCreatedAt, examRecords]) => {
            const scores = examRecords.map((examRecord) => examRecord.score);
            const maxScore = Math.max(...scores);
            myExamScores.push({
              _id: strCreatedAt.replace(/^\d{4}\//, ''), // yyyy/MM/dd -> MM/dd
              score: maxScore,
            });
          },
        );

        setExamScores(myExamScores);
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

    queryExamRecordOverview();
  }, [dispatch, toast, examId]);

  return (
    <Container maxW="container.xl" mt="3">
      {exam && (
        <>
          <PageHeading title={exam.topic} />
          <pre>{exam.description}</pre>
          <Divider my={5} />

          <Heading as="h3" size="lg" mt={5} style={{ textAlign: 'center' }}>
            {startDate} ~ Today
            <Tooltip label="若該日期有多筆分數，則以最高分為主">
              <InfoIcon ml={2} />
            </Tooltip>
          </Heading>

          <LineChart
            width={1200}
            height={300}
            data={examScores}
            margin={{
              top: 10,
            }}
          >
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis
              dataKey="_id"
              tick={{ fill: chartFontColor }}
              tickLine={{ stroke: chartFontColor }}
            />
            <YAxis
              tick={{ fill: chartFontColor }}
              tickLine={{ stroke: chartFontColor }}
              domain={[0, 10]}
            >
              <Label value="Score" fill={chartFontColor} />
            </YAxis>
            <ReChartTooltip labelStyle={{ color: 'black' }} />
            <Legend />
            <Line
              type="monotone"
              dataKey="score"
              stroke="#8884d8"
              activeDot={{ r: 6 }}
            />
          </LineChart>

          <Box borderWidth="1px" borderRadius="lg" p="2" mt="10">
            <Box>
              <Heading as="h3" size="lg" style={{ textAlign: 'center' }}>
                成績紀錄
              </Heading>
            </Box>
            <Box>
              <ExamRecordList examId={exam._id!} />
            </Box>
          </Box>

          <Box borderWidth="1px" borderRadius="lg" p="2" mt="10">
            <Box>
              <Heading as="h3" size="lg" style={{ textAlign: 'center' }}>
                答錯次數最多的問題
              </Heading>
            </Box>
            <Box>
              <WrongQuestionList questions={wrongQuestions} />
            </Box>
          </Box>
        </>
      )}
    </Container>
  );
}

export default ExamRecordOverview;
