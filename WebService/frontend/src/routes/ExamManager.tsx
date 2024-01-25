import { Container, Divider, ListItem, UnorderedList } from '@chakra-ui/react';
import ExamForm from '../components/exam/ExamForm';
import ExamList from '../components/exam/ExamList';
import PageHeading from '../components/PageHeading';
import ShowText from '../components/ShowText';

function ExamManager() {
  return (
    <Container maxW="container.xl" mt="3">
      <PageHeading title="測驗管理">
        <ShowText fontSize="lg">此頁可以新增、修改、刪除測驗</ShowText>
        <UnorderedList>
          <ListItem>
            <ShowText fontSize="lg">
              新增測驗 -
              輸入名稱(Topic)、描述(Descripton)後，點擊[新增]，即可增加自己的測驗
            </ShowText>
          </ListItem>
          <UnorderedList>
            <ListItem>
              <ShowText fontSize="lg">
                新增的測驗沒有題目，可在列表中點擊該測驗的[題目]按鍵，前往題目管理功能去新增題目
              </ShowText>
            </ListItem>
            <ListItem>
              <ShowText fontSize="lg">
                擁有題目的測驗，才會在首頁顯示出來
              </ShowText>
            </ListItem>
          </UnorderedList>
          <ListItem>
            <ShowText fontSize="lg">
              修改測驗 - 點擊[修改]會顯示對話框，修改內容後點擊[Confirm]
            </ShowText>
          </ListItem>
          <ListItem>
            <ShowText fontSize="lg">
              管理題目 -
              點擊[題目]會前往題目管理功能，用來新增、修改、刪除該筆測驗的題目
            </ShowText>
          </ListItem>
          <ListItem>
            <ShowText fontSize="lg">
              刪除測驗 - 點擊[刪除]會顯示確認對話框，點擊[Yes]將會刪除測驗
            </ShowText>
          </ListItem>
        </UnorderedList>
      </PageHeading>
      <ExamForm />
      <Divider my="5" />
      <ExamList />
    </Container>
  );
}

export default ExamManager;
