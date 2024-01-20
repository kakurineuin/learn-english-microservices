import { ReactNode, useRef } from 'react';
import {
  useDisclosure,
  Button,
  AlertDialog,
  AlertDialogOverlay,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogBody,
  AlertDialogFooter,
  AlertDialogCloseButton,
} from '@chakra-ui/react';

type ButtonColorSchema =
  | 'whiteAlpha'
  | 'blackAlpha'
  | 'gray'
  | 'red'
  | 'orange'
  | 'yellow'
  | 'green'
  | 'teal'
  | 'blue'
  | 'cyan'
  | 'purple'
  | 'pink'
  | 'linkedin'
  | 'facebook'
  | 'messenger'
  | 'whatsapp'
  | 'twitter'
  | 'telegram';
type ButtonVariant = 'ghost' | 'outline' | 'solid' | 'link' | 'unstyled';
type ButtonSize = 'lg' | 'md' | 'sm' | 'xs';

type Props = {
  button: {
    colorScheme: ButtonColorSchema;
    variant: ButtonVariant;
    size: ButtonSize;
    text: string;
  };
  header: string;
  children: ReactNode;
  yesCallback: () => void;
};

function ConfirmDialog({ button, header, children, yesCallback }: Props) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const cancelRef = useRef<HTMLButtonElement>(null);

  return (
    <>
      <Button
        colorScheme={button.colorScheme}
        variant={button.variant}
        size={button.size}
        onClick={onOpen}
      >
        {button.text}
      </Button>
      <AlertDialog
        motionPreset="slideInBottom"
        leastDestructiveRef={cancelRef}
        onClose={onClose}
        isOpen={isOpen}
        isCentered
      >
        <AlertDialogOverlay />

        <AlertDialogContent>
          <AlertDialogHeader>{header}</AlertDialogHeader>
          <AlertDialogCloseButton />
          <AlertDialogBody>{children}</AlertDialogBody>
          <AlertDialogFooter>
            <Button
              colorScheme="red"
              variant="outline"
              onClick={() => {
                yesCallback();
                onClose();
              }}
            >
              Yes
            </Button>
            <Button variant="outline" ml={3} ref={cancelRef} onClick={onClose}>
              No
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}

export default ConfirmDialog;
