import { Anchor, Box } from "@hope-ui/solid"
import summarize from "../utils/summarize"
import PipelineList from "../pipelines/PipelineList"

export default (props) => {
  return (
    <Box bgColor="$neutral3" borderRadius="$md" p="$3" boxShadow="$md">
      <Anchor href={props.repository.url} ml="$1" mr="$1" display="inline-block" external>
        <Box as="p" fontSize="$xs">
          {summarize(props.repository.namespace, 60)}
        </Box>
        <Box as="p" fontWeight="$semibold" fontSize="$base">
          {summarize(props.repository.name, 50)}
        </Box>
      </Anchor>
      <Box mt="$2">
        <PipelineList repositoryID={props.repository.id} />
      </Box>
    </Box>
  )
}
