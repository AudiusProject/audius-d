import { useState, DragEvent } from "react";
import { useAudiusSdk } from "../providers/AudiusSdkProvider";
import type { AudiusSdk as AudiusSdkType } from "@audius/sdk/dist/sdk/sdk.d.ts";
import { useAudiusLibs } from "../providers/AudiusLibsProvider";
import type { AudiusLibs } from "@audius/sdk/dist/WebAudiusLibs.d.ts";

import { DOMParser } from "linkedom";
import { Genre, UploadTrackRequest } from "@audius/sdk";

const fetchResource = async (url: string, filename: string) => {
  const res = await fetch(url);
  if (!res.ok) {
    throw new Error(
      `HTTP error when fetching clipper.jpg. Status: ${res.status}`,
    );
  }
  const blob = await res.blob();
  const mimeType = res.headers.get("Content-Type");
  return new File([blob], filename, { type: mimeType! });
};

const processXml = async (document: any, audiusSdk: AudiusSdkType) => {
  // todo remove this and upload images and tracks from xml without hardcoding
  const [clipperImg, snareAudio] = await Promise.all([
    fetchResource("/ddex-examples/clipper.jpg", "todo_img_name"),
    fetchResource("/ddex-examples/snare.wav", "todo_audio_name"),
  ]);

  // extract SoundRecording
  const trackNodes = queryAll(document, "SoundRecording", "track");

  for (const trackNode of Array.from(trackNodes)) {
    const tt = {
      title: firstValue(trackNode, "TitleText", "trackTitle"),

      // todo: need to normalize genre
      // genre: firstValue(trackNode, "Genre", "trackGenre"),
      genre: "Metal" as Genre,

      // todo: need to parse release date if present
      releaseDate: new Date(
        firstValue(trackNode, "OriginalReleaseDate", "originalReleaseDate"),
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
    const artistName = firstValue(trackNode, "ArtistName", "artistName");
    const { data: users } = await audiusSdk.users.searchUsers({
      query: artistName,
    });
    const userId = users[0].id;
    const uploadTrackRequest: UploadTrackRequest = {
      userId: userId,
      coverArtFile: clipperImg,
      metadata: tt,
      onProgress: (progress) => console.log("Progress:", progress),
      trackFile: snareAudio,
    };
    console.log(uploadTrackRequest);
    console.log("uploading track...");
    const result = await audiusSdk.tracks.uploadTrack(uploadTrackRequest);
    console.log(result);
  }

  // extract Release
  for (const releaseNode of queryAll(document, "Release", "release")) {
    //todo
  }
};

const queryAll = (node: any, ...fields: string[]) => {
  for (const field of fields) {
    const hits = node.querySelectorAll(field);
    if (hits.length) return Array.from(hits);
  }
  return [];
};

const firstValue = (node: any, ...fields: string[]) => {
  for (const field of fields) {
    const hit = node.querySelector(field);
    if (hit) return hit.textContent.trim();
  }
};

const validXmlFile = (file: File) => {
  return file.type === "text/xml" && file.name.endsWith(".xml");
};

const AudiusLogin = ({ audiusLibs }: { audiusLibs: AudiusLibs }) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loginError, setLoginError] = useState<string | null>(null);
  const [loginLoading, setLoginLoading] = useState(false);

  const handleLogin = async (e) => {
    setLoginLoading(true);
    e.preventDefault();
    try {
      const { error } = await audiusLibs.Account!.login(email, password);
      if (error) {
        throw new Error(error as string);
      }
    } catch (error) {
      setLoginError((error as Error).message);
    } finally {
      setLoginLoading(false);
    }
  };

  return (
    <div className="bg-gray-200 flex justify-center items-center h-screen">
      <div className="bg-white p-6 rounded shadow-md w-96 mx-auto">
        <form className="space-y-4">
          <h2 className="text-center text-2xl font-bold">Login to Audius</h2>
          <div>
            <label
              htmlFor="email"
              className="block mb-2 text-sm font-medium text-gray-700"
            >
              Email:
            </label>
            <input
              type="email"
              id="email"
              name="email"
              required
              className="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-blue-500"
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>
          <div>
            <label
              htmlFor="password"
              className="block mb-2 text-sm font-medium text-gray-700"
            >
              Password:
            </label>
            <input
              type="password"
              id="password"
              name="password"
              required
              className="w-full p-2 border border-gray-300 rounded focus:outline-none focus:border-blue-500"
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>
          <button
            type="submit"
            className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600"
            onClick={handleLogin}
            disabled={!email || !password || loginLoading}
          >
            {loginLoading ? "Logging in..." : "Login"}
          </button>
        </form>
        {loginError && <div className="text-red-500">{loginError}</div>}
      </div>
    </div>
  );
};

const XmlImporter = () => {
  const { audiusSdk } = useAudiusSdk();
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [uploadSucceeded, setUploadSucceeded] = useState(false);

  const handleDragIn = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handleDragOut = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  };

  const handleDragOver = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.dataTransfer.items && e.dataTransfer.items.length > 0) {
      setIsDragging(true);
    }
  };

  const handleDrop = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      handleFileChange(e.dataTransfer.files[0]); // reuse the file change handler
      e.dataTransfer.clearData();
    }
  };

  const clearSelection = () => {
    setSelectedFile(null);
  };

  const handleFileChange = (file: File) => {
    if (!validXmlFile(file)) {
      alert("Please upload an XML file.");
      return;
    }
    setSelectedFile(file);
    setUploadError(null);
  };

  const handleUpload = () => {
    if (!selectedFile) {
      alert("Please select a file first!");
      return;
    }

    if (!validXmlFile(selectedFile)) {
      alert("Please upload an XML file.");
      return;
    }

    readXml(selectedFile, audiusSdk!);
  };

  const readXml = (file: File, audiusSdk: AudiusSdkType) => {
    setIsUploading(true);
    const reader = new FileReader();
    reader.onload = async (event: ProgressEvent<FileReader>) => {
      const xmlText = event.target?.result;
      if (xmlText) {
        try {
          const document = new DOMParser().parseFromString(
            xmlText as string,
            "text/xml",
          );
          await processXml(document, audiusSdk);
          setUploadSucceeded(true);
        } catch (error) {
          setUploadError("Error processing xml");
          console.error("Error processing xml:", error);
        }
      } else {
        setUploadError("Error reading file");
      }
      setIsUploading(false);
    };
    reader.onerror = (error) => {
      setUploadError("Error reading file");
      console.error("Error reading file:", error);
      setIsUploading(false);
    };
    reader.readAsText(file);
  };

  return (
    <>
      <label
        className={`flex justify-center h-32 px-4 transition bg-white border-2 border-dashed rounded-md appearance-none cursor-pointer hover:border-gray-400 focus:outline-none ${
          isDragging ? "border-gray-400" : "border-gray-300 "
        }`}
        onDragEnter={handleDragIn}
        onDragLeave={handleDragOut}
        onDragOver={handleDragOver}
        onDrop={handleDrop}
      >
        <span className="flex items-center space-x-2">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="w-6 h-6 text-gray-600"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth="2"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
            />
          </svg>
          <span className="font-medium text-gray-600">
            {"Drop files to upload, or "}
            <span className="text-blue-600 underline">browse</span>
          </span>
        </span>
        <input
          type="file"
          name="file_upload"
          accept="text/xml,application/xml"
          className="hidden"
          onChange={(e) => handleFileChange(e.target.files[0])}
        />
      </label>
      {selectedFile && (
        <>
          <div>Selected file:</div>
          <div className="flex space-x-4">
            <div>{selectedFile.name}</div>
            <button
              className="text-xs w-8 p-1 bg-red-500 text-white rounded hover:bg-red-600 focus:outline-none"
              onClick={clearSelection}
            >
              x
            </button>
          </div>
          <button
            className="btn btn-blue"
            onClick={handleUpload}
            disabled={isUploading}
          >
            {isUploading ? "Uploading..." : "Upload"}
          </button>
          {uploadError && (
            <div className="text-red-500">
              {("Error uploading file: ", uploadError)}
            </div>
          )}
          {uploadSucceeded && (
            <div className="text-green-500">Upload success!</div>
          )}
        </>
      )}
    </>
  );
};

export const Ddex = () => {
  const { audiusLibs } = useAudiusLibs();
  return (
    <div className="flex flex-col space-y-4">
      {!audiusLibs || !audiusLibs.Account ? (
        "loading..."
      ) : !audiusLibs.Account.getCurrentUser() ? (
        <AudiusLogin audiusLibs={audiusLibs} />
      ) : (
        <>
          <div>{`Logged in as @${
            audiusLibs.Account.getCurrentUser().handle
          }`}</div>
          <XmlImporter />
        </>
      )}
    </div>
  );
};
