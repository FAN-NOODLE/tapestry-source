/**
 * el-date-picker 时间禁用的工厂函数
 *
 */
export const handleDisabledDate = (
  disabledRanges: Array<any> = [],
  selectStart?: any
) => {
  return (val: any) => {
    const now = new Date().setHours(0, 0, 0, 0) / 1000;
    const current = val.setHours(0, 0, 0, 0) / 1000;
    if (!disabledRanges.length) {
      return current < now;
    }
    if (selectStart) {
      const selectStartTime = selectStart.setHours(0, 0, 0, 0) / 1000;
      if (selectStartTime < disabledRanges[0].distribution_start) {
        return current >= disabledRanges[0].distribution_start || current < now;
      }
      if (
        selectStartTime >
        disabledRanges[disabledRanges.length - 1].distribution_end
      ) {
        return (
          current <=
            disabledRanges[disabledRanges.length - 1].distribution_end ||
          current < now
        );
      }

      for (let i = 0; i < disabledRanges.length - 1; i++) {
        if (
          selectStartTime > disabledRanges[i].distribution_end &&
          selectStartTime < disabledRanges[i + 1].distribution_start
        ) {
          return (
            current <= disabledRanges[i].distribution_end ||
            current >= disabledRanges[i + 1].distribution_start ||
            current < now
          );
        }
      }
    } else {
      return (
        disabledRanges.some(item => {
          return (
            current >= item.distribution_start &&
            current <= item.distribution_end
          );
        }) || current < now
      );
    }
  };
};
