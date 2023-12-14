import React, { useState } from 'react';
import { useAudiusSdk } from "../providers/AudiusSdkProvider";
import "dotenv/config";
import type { AudiusSdk as AudiusSdkType } from "@audius/sdk/dist/sdk/sdk.d.ts";

import { DOMParser } from "linkedom";
import fs from "fs";
import {
  Genre,
  UploadTrackRequest,
} from "@audius/sdk";

const processXml = async (document: any, audiusSdk: AudiusSdkType) => {
  try {
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
      //todo
    }
  } catch (error) {
    console.error("Error processing XML:", error);
  }
}

const readXml = (file: Blob, audiusSdk: AudiusSdkType) => {
  const reader = new FileReader();
  reader.onload = async (event) => {
    const xmlText = event.target.result;
    console.log(xmlText);
    const document = new DOMParser().parseFromString(
      xmlText,
      "text/xml",
    );
    console.log(document);
    await processXml(document, audiusSdk);
  };
  reader.onerror = (error) => {
    //todo
    console.error('Error reading file:', error);
  }
  reader.readAsText(file);
}

const queryAll = (node: any, ...fields: string[]) => {
  for (const field of fields) {
    const hits = node.querySelectorAll(field);
    if (hits.length) return Array.from(hits);
  }
  return [];
}

const firstValue = (node: any, ...fields: string[]) => {
  for (const field of fields) {
    const hit = node.querySelector(field);
    if (hit) return hit.textContent.trim();
  }
}

const validXmlFile = (file: Blob) => {
  return file.type === 'text/xml' && file.name.endsWith('.xml');
}

export const XmlImporter = () => {
  const { audiusSdk } = useAudiusSdk();
  const [selectedFile, setSelectedFile] = useState(null);
  const [isDragging, setIsDragging] = useState(false);

  const handleDragIn = (e) => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handleDragOut = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  };

  const handleDragOver = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.dataTransfer.items && e.dataTransfer.items.length > 0) {
      setIsDragging(true);
    }
  };

  const handleDrop = (e) => {
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
  }

  const handleFileChange = (file: Blob) => {
    if (!validXmlFile(file)) {
      alert('Please upload an XML file.');
      return;
    }
    setSelectedFile(file);
  };

  const handleUpload = () => {
    if (!selectedFile) {
      alert('Please select a file first!');
      return;
    }

    if (!validXmlFile(selectedFile)) {
      alert('Please upload an XML file.');
      return;
    }

    readXml(selectedFile, audiusSdk!);
  };

  return (
    <div className="flex flex-col space-y-4">
      <label
        className={`flex justify-center h-32 px-4 transition bg-white border-2 border-dashed rounded-md appearance-none cursor-pointer hover:border-gray-400 focus:outline-none ${isDragging ? 'border-gray-400' : 'border-gray-300 '}`}
        onDragEnter={handleDragIn} 
        onDragLeave={handleDragOut} 
        onDragOver={handleDragOver} 
        onDrop={handleDrop}
      >
          <span className="flex items-center space-x-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-gray-600" fill="none" viewBox="0 0 24 24"
                  stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round"
                      d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
              </svg>
              <span className="font-medium text-gray-600">
                {'Drop files to upload, or '}
                <span className="text-blue-600 underline">browse</span>
              </span>
          </span>
          <input type="file" name="file_upload" accept="text/xml,application/xml" class="hidden" onChange={(e) => handleFileChange(e.target.files[0])} />
      </label>
      {selectedFile && (
        <div>
          <div>
            Selected file:
          </div>
          <div className="flex space-x-4">
            <div>
              {selectedFile.name}
            </div>
            <button className="text-xs w-8 p-1 bg-red-500 text-white rounded hover:bg-red-600 focus:outline-none" onClick={clearSelection}>x</button>
          </div>
        </div>
      )}
      <button className="btn btn-blue" onClick={handleUpload}>Upload</button>
    </div>
  );
}
