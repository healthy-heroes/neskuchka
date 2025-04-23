import { motion } from "framer-motion";
import {
  Heading,
  Flex,
  Text,
  Button,
  ButtonGroup,
  View,
} from "@adobe/react-spectrum";

import { pageProps } from "../../constants";

import styles from "./Intro.module.css";

export function Intro() {
  const motionConfig = {
    container: {
      hidden: { opacity: 0 },
      show: {
        opacity: 1,
        transition: {
          staggerChildren: 0.2,
        },
      },
    },
    item: {
      hidden: { opacity: 0, y: 20 },
      show: { opacity: 1, y: 0 },
    },

    button: {
      whileHover: { y: -2 },
      whileTap: { scale: 0.95 },
    },
  };

  const handleStart = () => {
    console.log("Start");
  };

  const handleLearnMore = () => {
    console.log("Learn More");
  };

  return (
    <View
      paddingX={pageProps.paddingX}
      paddingY="size-600"
      backgroundColor="indigo-400"
    >
      <motion.div
        variants={motionConfig.container}
        initial="hidden"
        animate="show"
      >
        <Flex
          direction="column"
          gap="size-400"
          alignItems="center"
          justifyContent="center"
          maxWidth={pageProps.maxWidth}
          marginX="auto"
        >
          <motion.div variants={motionConfig.item}>
            <Heading
              level={1}
              marginBottom="size-1"
              marginTop="size-200"
              UNSAFE_className={styles.title}
            >
              Откройте для себя мир Нескучного Спорта
            </Heading>
          </motion.div>

          <motion.div variants={motionConfig.item}>
            <Text UNSAFE_className={styles.subtitle}>
              Бесплатные тренировки + анализ твоей физической активности!
            </Text>
          </motion.div>

          <motion.div variants={motionConfig.item}>
            <ButtonGroup marginTop="size-400">
              <Button
                variant="accent"
                UNSAFE_className={styles.startButton}
                onPress={handleStart}
              >
                Начать
              </Button>
              <Button
                variant="primary"
                staticColor="white"
                UNSAFE_className={styles.featuresButton}
                onPress={handleLearnMore}
              >
                Узнать больше
              </Button>
            </ButtonGroup>
          </motion.div>
        </Flex>
      </motion.div>
    </View>
  );
}
