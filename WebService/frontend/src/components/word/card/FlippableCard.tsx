import {
  Box,
  Heading,
  Text,
  Card,
  CardHeader,
  CardBody,
  ListItem,
  UnorderedList,
} from '@chakra-ui/react';
import styles from './FlippableCard.module.css';
import { WordMeaning } from '../../../models/WordMeaning';
import AudioButton from '../../AudioButton';

type Props = {
  isFlip: boolean;
  wordMeaning: WordMeaning;
};

function FlippableCard({ isFlip, wordMeaning: { word, examples } }: Props) {
  const classNames = [styles.card];

  if (isFlip) {
    classNames.push(styles.flip);
  }

  const exampleComponents = examples.map(
    ({ pattern, examples: exampleArray }) => {
      const components = exampleArray.map(({ audioUrl, text: sentence }) => (
        <ListItem mt="2" key={sentence}>
          <AudioButton
            colorScheme="teal"
            variant="outline"
            size="sm"
            mr="2"
            audioUrl={audioUrl}
          />
          {sentence}
        </ListItem>
      ));

      if (pattern) {
        return (
          <Box my="7" key={pattern}>
            <Text as="b" fontSize="xl">
              {pattern}
            </Text>
            <UnorderedList>{components}</UnorderedList>
          </Box>
        );
      }

      return components[0];
    }
  );

  return (
    <div className={styles['card-container']}>
      <div className={classNames.join(' ')}>
        {/* 卡片反面 */}
        <div className={styles['card-back']}>
          <Card variant="outline" w="100%" h="100%">
            <CardHeader alignSelf="center">
              <Heading as="h1" size="4xl">
                {word}
              </Heading>
            </CardHeader>

            <CardBody overflowY="auto">
              <UnorderedList mt="5">{exampleComponents}</UnorderedList>
            </CardBody>
          </Card>
        </div>

        {/* 卡片正面 */}
        <div className={styles['card-front']}>
          <div className={styles['card-front-pattern']} />
        </div>
      </div>
    </div>
  );
}

export default FlippableCard;
