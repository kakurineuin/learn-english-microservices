import {
  Container,
  Divider,
  ListItem,
  Text,
  UnorderedList,
} from '@chakra-ui/react';
import { useLocation } from 'react-router-dom';
import QuestionForm from '../components/exam/QuestionForm';
import QuestionList from '../components/exam/QuestionList';
import { Exam } from '../models/Exam';
import PageHeading from '../components/PageHeading';
import ShowText from '../components/ShowText';

type LocationState = {
  exam: Exam;
};

function QuestionManager() {
  const {
    state: { exam },
  }: { state: LocationState } = useLocation();

  return (
    <Container maxW="container.xl" mt="3">
      <PageHeading title={exam.topic}>
        <ShowText fontSize="lg">此頁可以新增、修改、刪除題目</ShowText>
        <UnorderedList>
          <ListItem>
            <ShowText fontSize="lg">
              新增題目 -
              輸入問題(Ask)、答案(Answer)，若點選[+]可增加要求輸入的答案，然後點擊[新增]
            </ShowText>
          </ListItem>
          <ListItem>
            <ShowText fontSize="lg">
              修改題目 - 點擊[修改]會顯示對話框，修改內容後點擊[Confirm]
            </ShowText>
          </ListItem>
          <ListItem>
            <ShowText fontSize="lg">
              刪除題目 - 點擊[刪除]會顯示確認對話框，點擊[Yes]將會刪除題目
            </ShowText>
          </ListItem>
        </UnorderedList>
      </PageHeading>
      <Text>{exam.description}</Text>
      <Divider my="5" />
      <QuestionForm examId={exam._id!} />
      <Divider my="5" />
      <QuestionList examId={exam._id!} />
    </Container>
  );
}

export default QuestionManager;
