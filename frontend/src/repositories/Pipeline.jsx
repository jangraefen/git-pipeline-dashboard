import { Box, HStack, Skeleton, Tag, TagLeftIcon, TagLabel, Anchor } from "@hope-ui/solid"
import { BiSolidArrowFromLeft, BiRegularCalendar, BiRegularGitBranch } from "solid-icons/bi"
import { createEffect, createReaction, createResource, createSignal, onCleanup, Show } from "solid-js"
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

const pipelineFetcher = async ({ repositoryID, repositorySource }) =>
  (await fetch(`/repositories/${repositorySource}/${repositoryID}`)).json()

export default (props) => {
  const [pipeline, { refetch: refetchPipeline }] = createResource(() => ({ ...props }), pipelineFetcher)
  const interval = setInterval(refetchPipeline, 1 * 60 * 1000)
  onCleanup(() => clearInterval(interval))

  const [boxClass, setBoxClass] = createSignal("")
  createEffect(() => {
    if (pipeline() != null) {
      setBoxClass(stateToStyle(pipeline().state))
    }
  })

  const [initiallyLoaded, setInitiallyLoaded] = createSignal(false)
  const loadTracker = createReaction(() => setInitiallyLoaded(true))
  loadTracker(pipeline)

  return (
    <Box>
      <Skeleton loaded={initiallyLoaded}>
        <Show when={pipeline.error == null && pipeline() != null && !pipeline().status}>
          <Box class={boxClass() + " " + (pipeline.loading ? styles.updating : "")} borderRadius="$md" p="$3">
            <Anchor href={pipeline().url} external>
              <Box mt="$2">
                <Box fontSize="$xs">{summarize(pipeline().commitMessage, 60)}</Box>
                <Box fontSize="$xs">by {summarize(pipeline().commitAuthor, 60)}</Box>
              </Box>
            </Anchor>
            <HStack spacing="$1" mt="$4">
              <Tag colorSchema="neutral">
                <TagLeftIcon as={BiRegularCalendar} />
                <TagLabel fontSize="$2xs">{pipeline().time}</TagLabel>
              </Tag>
              <Tag colorSchema="neutral">
                <TagLeftIcon as={BiRegularGitBranch} />
                <TagLabel fontSize="$2xs">{summarize(pipeline().ref, 20)}</TagLabel>
              </Tag>
              <Tag colorSchema="neutral">
                <TagLeftIcon as={BiSolidArrowFromLeft} />
                <TagLabel fontSize="$2xs">{pipeline().state}</TagLabel>
              </Tag>
            </HStack>
          </Box>
        </Show>
        <Show when={pipeline.error == null && pipeline() != null && !!pipeline().status}>
          <Box bg="$neutral7" borderRadius="$md" p="$3" fontStyle="italic" fontSize="$xs">
            <Box mt="$2">
              <Box fontSize="$xs">Error while fetching pipeline:</Box>
              <Box fontSize="$xs">
                {pipeline().state} - {pipeline().title}
              </Box>
              <Box fontSize="$xs">{pipeline().detail}</Box>
            </Box>
          </Box>
        </Show>
        <Show when={pipeline.error == null && pipeline() == null}>
          <Box bg="$neutral7" borderRadius="$md" p="$3" fontStyle="italic" fontSize="$xs">
            No pipelines detected.
          </Box>
        </Show>
        <Show when={pipeline.error != null}>
          <Box bg="$neutral7" borderRadius="$md" p="$3" fontStyle="italic" fontSize="$xs">
            <Box mt="$2">
              <Box fontSize="$xs">Error fetching pipeline information:</Box>
              <Box fontSize="$xs">by {pipeline.error}</Box>
            </Box>
          </Box>
        </Show>
      </Skeleton>
    </Box>
  )
}
