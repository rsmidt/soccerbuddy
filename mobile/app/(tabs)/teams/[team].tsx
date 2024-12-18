/**
 * This is needed so Expo thinks this routes exists even though we're intercepting
 * any route requests manually in the _layout file.
 * This is another ugly hack and shows how minimal useful Expo Router is.
 */
export default function TeamIndex() {
  return null;
}
