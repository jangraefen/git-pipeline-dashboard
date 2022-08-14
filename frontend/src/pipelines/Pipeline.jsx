import { Box, HStack, Tag, TagLeftIcon, TagLabel, Anchor, Text } from "@hope-ui/solid"
import { BiSolidArrowFromLeft, BiRegularCalendar, BiRegularGitBranch } from "solid-icons/bi"
import { createEffect, createSignal } from "solid-js"
import summarize from "../utils/summarize"
import styles from "./Pipeline.module.css"

const stateToStyle = (state) => {
  if (state === "success") {
    return styles.success
  } else if (state === "failed") {
    return styles.failed
  } else if (state === "running") {
    return styles.running
  }

  return styles.unknown
}

export default (props) => {
  const [boxClass, setBoxClass] = createSignal("")
  createEffect(() => {
    if (props.pipeline != null) {
      setBoxClass(stateToStyle(props.pipeline.state))
    }
  })

  return (
    <Box class={boxClass() + " " + (props.pipeline.loading ? styles.updating : "")} borderRadius="$md" p="$3">
      <Text>{props.pipeline.name}</Text>
      <Anchor href={props.pipeline.url} external>
        <Box mt="$2">
          <Box fontSize="$xs">{summarize(props.pipeline.commitMessage, 60)}</Box>
          <Box fontSize="$xs">by {summarize(props.pipeline.commitAuthor, 60)}</Box>
        </Box>
      </Anchor>
      <HStack spacing="$1" mt="$4">
        <Tag colorSchema="neutral">
          <TagLeftIcon as={BiRegularCalendar} />
          <TagLabel fontSize="$2xs">{props.pipeline.time}</TagLabel>
        </Tag>
        <Tag colorSchema="neutral">
          <TagLeftIcon as={BiRegularGitBranch} />
          <TagLabel fontSize="$2xs">{summarize(props.pipeline.ref, 20)}</TagLabel>
        </Tag>
        <Tag colorSchema="neutral">
          <TagLeftIcon as={BiSolidArrowFromLeft} />
          <TagLabel fontSize="$2xs">{props.pipeline.state}</TagLabel>
        </Tag>
      </HStack>
    </Box>
  )
}
