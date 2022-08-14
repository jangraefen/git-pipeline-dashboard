export default (text, maxLength) => {
  if (text == null || text === "") {
    return ""
  }

  const length = Math.min(text.length, maxLength - 3)

  let summarizedText = text.substring(0, length)
  if (length < text.length) {
    summarizedText += "..."
  }

  return summarizedText
}
