import { Text } from "@adobe/react-spectrum";

import { LogoIcon } from "./LogoIcon";
import styles from "./Logo.module.css";

export function Logo() {
  return (
    <Text UNSAFE_className={styles.logo}>
      <LogoIcon size="S" position="relative" top="size-25" />
      &nbsp;Neskuchka
    </Text>
  );
}
