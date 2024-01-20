import { Container } from '@chakra-ui/react';
import UserHistoryList from '../components/user/UserHistoryList';
import PageHeading from '../components/PageHeading';
import ShowText from '../components/ShowText';

function UserHistoryManager() {
  return (
    <Container maxW="container.xl" mt="3">
      <PageHeading title="使用者歷史紀錄">
        <ShowText fontSize="lg">此頁可以查看使用者過去的操作紀錄</ShowText>
      </PageHeading>
      <UserHistoryList />
    </Container>
  );
}

export default UserHistoryManager;
