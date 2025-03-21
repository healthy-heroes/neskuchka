import { Flex, Text, Button, View, ButtonGroup } from "@adobe/react-spectrum";

import { pageProps } from "../../pages/constants";
import { Logo } from "./Logo/Logo";

import styles from "./NavBar.module.css";

export function NavBar() {
  const handleSignIn = () => {
    console.log("Sign In clicked");
  };

  const handleSignUp = () => {
    console.log("Sign Up clicked");
  };

  const menuItems = ["Тренировки", "О нас"];

  return (
    <View backgroundColor="gray-50" paddingX="size=100" paddingY="size-200">
      <Flex
        direction="row"
        alignItems="center"
        justifyContent="space-between"
        maxWidth={pageProps.maxWidth}
        marginX="auto"
        width="100%"
      >
        <Logo />

        <Flex gap="size-400" alignItems="center">
          <Flex gap="size-300" alignItems="center">
            {menuItems.map((item) => (
              <Text key={item} UNSAFE_className={styles.menuLink}>
                {item}
              </Text>
            ))}
          </Flex>

          <ButtonGroup>
            <Button variant="primary" onPress={handleSignIn}>
              Войти
            </Button>
            <Button variant="accent" onPress={handleSignUp}>
              Зарегистрироваться
            </Button>
          </ButtonGroup>
        </Flex>
      </Flex>
    </View>
  );
}
