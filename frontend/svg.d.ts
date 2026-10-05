declare module "*.svg" {
  import React from "react";
  const Component: React.FunctionComponent<React.SVGProps<SVGSVGElement>>;
  export default Component;
}

declare module "*.svg?url" {
  const src: string;
  export default src;
}