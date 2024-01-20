import { Container, ListItem, UnorderedList } from '@chakra-ui/react';
import PageHeading from '../components/PageHeading';
import ShowText from '../components/ShowText';
import WordForm from '../components/word/WordForm';

function WordManager() {
  return (
    <Container maxW="container.xl" mt="3">
      <PageHeading title="查詢單字">
        <ShowText fontSize="lg">
          此頁可以查詢英文單字的解釋，資料來自朗文線上字典(https://www.ldoceonline.com/)
        </ShowText>
        <UnorderedList>
          <ListItem>
            <ShowText fontSize="lg">
              查詢到的解釋，可以點擊[Favorite]加入到最愛的單字解釋，再點擊一次即可取消
            </ShowText>
          </ListItem>
        </UnorderedList>
      </PageHeading>
      <WordForm />
    </Container>
  );
}

export default WordManager;
