import { Accordion, AccordionButton, AccordionIcon, AccordionItem, AccordionPanel, Stack, Text } from "@hope-ui/solid"
import { For } from "solid-js"
import Repository from "./Repository"

import styles from "./RepositoryGroup.module.css"

export default (props) => {
  return (
    <Accordion defaultIndex={0} w="100%">
      <AccordionItem class={styles.repositoryGroup}>
        <AccordionButton fontWeight="$semibold">
          <AccordionIcon />
          <Text flex={1} textAlign="start">
            {props.title}
          </Text>
        </AccordionButton>
        <AccordionPanel>
          <Stack spacing="$5" flexWrap="wrap">
            <For each={props.repositories}>{(repository) => <Repository repository={repository} />}</For>
          </Stack>
        </AccordionPanel>
      </AccordionItem>
    </Accordion>
  )
}
