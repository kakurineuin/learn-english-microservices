import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Badge, Button, Box, Center, Flex, Spacer } from '@chakra-ui/react';
import CardWithScale from '../CardWithScale';
import ShowText from '../ShowText';
import { useAppDispatch, useAppSelector } from '../../store/hooks';
import { askFormActions } from '../../store/slices/askFormSlice';

type Props = {
  examId: string;
  topic: string;
  description: string;
  isPublic: boolean;
  questionCount: number;
  recordCount: number;
};

function ExamInfoCard({
  examId,
  topic,
  description,
  isPublic,
  questionCount,
  recordCount,
}: Props) {
  const user = useAppSelector((state) => state.session.user);
  const navigate = useNavigate();
  const dispatch = useAppDispatch();

  const goStartExamClickHandler = () => {
    // 避免重複計分，需要先清除 askFormSlice 的 questionResult
    dispatch(askFormActions.clearQuestionResult());
    navigate(`/restricted/exam/${examId}/start`);
  };
  const goExamRecordOverviewClickHandler = () => {
    navigate(`/restricted/exam/${examId}/record/overview`);
  };

  const [isShowDescription, setIsShowDescription] = useState(false);
  const showDescriptionOnAnimationComplete = () => {
    // 動畫結束後，再延遲一下下，感覺上會好一點
    setTimeout(() => setIsShowDescription(true), 500);
  };

  return (
    <CardWithScale
      onAnimationComplete={() => showDescriptionOnAnimationComplete()}
    >
      <Box
        boxShadow="xl"
        maxW="sm"
        borderWidth="1px"
        borderRadius="lg"
        overflow="hidden"
      >
        <Box p="6">
          <Flex>
            <Box mt="1" fontWeight="semibold" as="h4" lineHeight="tight">
              {topic}
            </Box>
            <Spacer />
            {isPublic && (
              <Center>
                <Badge variant="outline" colorScheme="green">
                  Public
                </Badge>
              </Center>
            )}
          </Flex>

          <Box>{isShowDescription && <ShowText>{description}</ShowText>}</Box>

          <Box
            color="gray.500"
            fontWeight="semibold"
            letterSpacing="wide"
            fontSize="xs"
            textTransform="uppercase"
            ml="2"
          >
            {questionCount} questions
          </Box>

          <Box
            color="gray.500"
            fontWeight="semibold"
            letterSpacing="wide"
            fontSize="xs"
            textTransform="uppercase"
            ml="2"
          >
            {recordCount} records
          </Box>

          <Box mt="4">
            {questionCount > 0 && (
              <Button
                mr="3"
                variant="outline"
                colorScheme="teal"
                size="sm"
                isDisabled={!user}
                onClick={goStartExamClickHandler}
              >
                開始測驗
              </Button>
            )}

            {recordCount > 0 && (
              <Button
                variant="outline"
                colorScheme="blue"
                size="sm"
                isDisabled={!user}
                onClick={goExamRecordOverviewClickHandler}
              >
                成績紀錄
              </Button>
            )}
          </Box>
        </Box>
      </Box>
    </CardWithScale>
  );
}

export default ExamInfoCard;
