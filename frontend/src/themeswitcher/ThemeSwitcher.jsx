import { HStack, Icon, Switch, useColorMode } from "@hope-ui/solid"
import { FaRegularMoon, FaRegularSun } from "solid-icons/fa"

export default () => {
  const { colorMode, toggleColorMode } = useColorMode()

  return (
    <HStack spacing="$2">
      <Icon as={FaRegularSun} />
      <Switch m="0" p="0" b="0" colorScheme="neutral" checked={colorMode() === "dark"} onChange={toggleColorMode} />
      <Icon as={FaRegularMoon} />
    </HStack>
  )
}
