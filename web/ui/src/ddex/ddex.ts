import "dotenv/config";
import type { AudiusSdk as AudiusSdkType } from "@audius/sdk/dist/sdk/sdk.d.ts";

import { DOMParser } from "linkedom";
import fs from "fs";
import {
  Genre,
  UploadTrackRequest,
} from "@audius/sdk";

export const importXml = async (fileName: string, audiusSdk: AudiusSdkType) => {
  // await await delay(1000);
  console.log(fileName);
  const xmlText = await fs.promises.readFile(fileName);
  const document = new DOMParser().parseFromString(
    xmlText.toString(),
    "text/xml",
  );

  // extract SoundRecording
  const trackNodes = queryAll(document, "SoundRecording", "track");

  for (const trackNode of Array.from(trackNodes)) {
    const tt = {
      title: useFirstValue(trackNode, "TitleText", "trackTitle"),

      // todo: need to normalize genre
      // genre: useFirstValue(trackNode, "Genre", "trackGenre"),
      genre: "Metal" as Genre,

      // todo: need to parse release date if present
      releaseDate: new Date(
        useFirstValue(trackNode, "OriginalReleaseDate", "originalReleaseDate"),
      ),
      // releaseDate: new Date(),

      isUnlisted: false,
      isPremium: false,
      fieldVisibility: {
        genre: true,
        mood: true,
        tags: true,
        share: true,
        play_count: true,
        remixes: true,
      },
      description: "",
      license: "Attribution ShareAlike CC BY-SA",
    };
    const artistName = useFirstValue(trackNode, "ArtistName", "artistName");
    const { data: users } = await audiusSdk.users.searchUsers({
      query: artistName,
    });
    const userId = users[0].id;
    const uploadTrackRequest: UploadTrackRequest = {
      userId: userId,
      coverArtFile: {
        buffer: await fs.promises.readFile("src/ddex/examples/clipper.jpg"),
        name: "todo_file_name",
      },
      metadata: tt,
      onProgress: (progress) => console.log("Progress:", progress),
      trackFile: {
        buffer: await fs.promises.readFile("src/ddex/examples/snare.wav"),
        name: "todo_track_file_name",
      },
    };
    console.log("uploading track...");
    const result = await audiusSdk.tracks.uploadTrack(uploadTrackRequest);
    console.log(result);
  }

  // extract Release
  for (const releaseNode of queryAll(document, "Release", "release")) {
  }
}

function queryAll(node: any, ...fields: string[]) {
  for (const field of fields) {
    const hits = node.querySelectorAll(field);
    if (hits.length) return Array.from(hits);
  }
  return [];
}

function useFirstValue(node: any, ...fields: string[]) {
  for (const field of fields) {
    const hit = node.querySelector(field);
    if (hit) return hit.textContent.trim();
  }
}

// ------------------------- main entry point stuff

// importXml("src/ddex/examples/DistroKid.xml");
// importXml("examples/Other.xml");
// importXml("20230901000527673_update/196834227253.xml");
