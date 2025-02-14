import { motion } from "framer-motion";
import {
  Heading,
  Content,
  Flex,
  Text,
  Button,
  ButtonGroup,
  View,
} from "@adobe/react-spectrum";

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
    <View paddingY="size-1000" backgroundColor="indigo-400">
      <Content>
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
            maxWidth={{ base: "100%", L: "1200px" }}
            marginX="auto"
          >
            <motion.div variants={motionConfig.item}>
              <Heading level={1} UNSAFE_className={styles.title}>
                Тренируйся с удовольствием
              </Heading>
            </motion.div>

            <motion.div variants={motionConfig.item}>
              <Text UNSAFE_className={styles.subtitle}>
                Тренировки для каждого. Тренируйся в своем темпе, получай
                удовольствие от результата.
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
      </Content>
    </View>
  );
}
