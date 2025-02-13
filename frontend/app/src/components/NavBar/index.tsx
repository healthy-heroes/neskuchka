import {
  Header,
  Flex,
  Text,
  Button,
  View,
  ButtonGroup,
} from "@adobe/react-spectrum";

import styles from "./NavBar.module.css";

export function NavBar() {
  const handleSignIn = () => {
    console.log("Sign In clicked");
  };

  const handleSignUp = () => {
    console.log("Sign Up clicked");
  };

  const menuItems = ["Workouts", "About"];

  return (
    <Header>
      <View paddingX="size=100" paddingY="size-100" backgroundColor="gray-50">
        <Flex
          direction="row"
          alignItems="center"
          justifyContent="space-between"
          maxWidth={{ base: "100%", L: "1200px" }}
          marginX="auto"
          width="100%"
        >
          <Text UNSAFE_className={styles.Logo}>Neskuchka</Text>

          <Flex gap="size-400" alignItems="center">
            <Flex gap="size-300" alignItems="center">
              {menuItems.map((item) => (
                <Text key={item} UNSAFE_className={styles.MenuLink}>
                  {item}
                </Text>
              ))}
            </Flex>

            <ButtonGroup>
              <Button variant="primary" onPress={handleSignIn}>
                Sign In
              </Button>
              <Button variant="accent" onPress={handleSignUp}>
                Sign Up
              </Button>
            </ButtonGroup>
          </Flex>
        </Flex>
      </View>
    </Header>
  );
}
