import { ethers } from "ethers";

export const utf8ToBytes32 = (utf8Str: string) => {
  return ethers.encodeBytes32String(utf8Str);
};

export const bytes32ToUtf8 = (bytes32Str: string) => {
  return ethers.decodeBytes32String(bytes32Str);
};
