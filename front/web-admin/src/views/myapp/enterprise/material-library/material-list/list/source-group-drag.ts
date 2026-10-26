export function shouldRequestSourceGroupMove(
  dropNode: any,
  dropType: string,
  previousSiblingIds: any[],
  currentSiblingIds: any[]
) {
  if (!dropNode || dropType === "none") {
    return false;
  }

  return (
    previousSiblingIds.length !== currentSiblingIds.length ||
    previousSiblingIds.some(
      (id: any, index: number) => id !== currentSiblingIds[index]
    )
  );
}
