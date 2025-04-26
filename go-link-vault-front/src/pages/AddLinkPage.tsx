import React, { useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

const AddLinkPage = () => {
  const [url, setUrl] = useState("");
  const [title, setTitle] = useState("");
  const [tags, setTags] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    try {
      const token = localStorage.getItem("token");
      const response = await axios.post(
        "http://localhost:8080/links",
        {
          url,
          title,
          tags: tags.split(",").map((t) => t.trim()),
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      navigate("/"); // Redirige vers la home après ajout
    } catch (err: any) {
      setError(err.response?.data?.message || "Something went wrong");
    }
  };

  return (
    <div className="max-w-md mx-auto mt-10">
      <h1 className="text-2xl font-bold mb-4 text-center">Add a Link</h1>
      {error && <p className="text-red-500 mb-4">{error}</p>}
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm">URL</label>
          <input
            type="url"
            className="w-full p-2 border rounded"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            required
          />
        </div>
        <div>
          <label className="block text-sm">Title</label>
          <input
            className="w-full p-2 border rounded"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
        </div>
        <div>
          <label className="block text-sm">Tags (comma-separated)</label>
          <input
            className="w-full p-2 border rounded"
            value={tags}
            onChange={(e) => setTags(e.target.value)}
          />
        </div>
        <button className="w-full bg-blue-500 text-white py-2 rounded" type="submit">
          Add Link
        </button>
      </form>
    </div>
  );
};

export default AddLinkPage;

