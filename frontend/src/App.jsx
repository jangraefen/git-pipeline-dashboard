import { Container, Flex, Heading, Spacer, VStack } from "@hope-ui/solid"
import RepositoryList from "./repositories/RepositoryList"
import ThemeSwitcher from "./themeswitcher/ThemeSwitcher"

function App() {
  return (
    <Container>
      <VStack spacing="$4">
        <Flex w="100%" mt="$2">
          <Heading>GIT Pipeline Dashboard</Heading>
          <Spacer />
          <ThemeSwitcher />
        </Flex>
        <RepositoryList />
      </VStack>
    </Container>
  )
}

export default App
