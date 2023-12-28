import { ReactNode } from 'react';
import { Center, Heading } from '@chakra-ui/react';
import InfoDialog from './InfoDialog';

type Props = {
  title: String;
  children?: ReactNode;
};

function PageHeading({ title, children }: Props) {
  return (
    <Center my="3">
      <Heading>{title}</Heading>
      {children && <InfoDialog>{children}</InfoDialog>}
    </Center>
  );
}

export default PageHeading;
