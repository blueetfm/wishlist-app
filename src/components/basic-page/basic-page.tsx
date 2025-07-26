import React from "react";
import Navbar from "../navbar/navbar";

export const BasicPage = ({children} : {children: React.ReactNode}) => {
    return (
        <>
        <Navbar />
        {children}
        </>
    )
}