import { initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";

const firebaseConfig = {
  apiKey: "AIzaSyDwwJV2NDbiSEKXqZOzI6vokZL_1VTmKu8",
  authDomain: "ticket-auth-58b02.firebaseapp.com",
  projectId: "ticket-auth-58b02",
  storageBucket: "ticket-auth-58b02.firebasestorage.app",
  messagingSenderId: "241066650036",
  appId: "1:241066650036:web:5b0b825a1209d396b72ae1",
  measurementId: "G-RZ10SFMBTC",
};

const app = initializeApp(firebaseConfig);
export const auth = getAuth(app);
