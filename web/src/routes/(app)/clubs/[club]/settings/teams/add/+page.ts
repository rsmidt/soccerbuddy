import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, data }) => {
  const { club } = await parent();
  return {
    club,
    ...data,
  };
};
