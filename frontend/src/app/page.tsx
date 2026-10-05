import Image from "next/image";
import WishlistItem from "@/components/wishlist-item/wishlist-item-detail";
import { AuthButton } from "@/components/auth-button/auth-button";
import { StylizedButton } from "@/components/ui/stylized-button";
import fujifilmImg from '../assets/fujifilm.webp';
import "./page.css";

export default function Page() {
  return (
      <main className="page">
        
        {/* Top Navigation Bar inside the rounded rectangle */}
        <header className="page__header">
          {/* Top Left: Logo */}
          <div className="page__logo">
            <Image src="/logo-name.png" alt="Logo" width={125} height={125}  />
          </div>
          {/* Top Right: Auth Button */}
          <div className="page__auth">
            <AuthButton />
          </div>
        </header>

        {/* Main Content Split */}
        <div className="page__content">
          
          {/* Left Column: WishlistItem Example */}
          <div className="page__left-column">
            <div className="page__rotate-wrapper">
              <WishlistItem 
                name="Fujifilm X100VI"
                description="been wanting this for ages..."
                price={1599}
                imageUrl={fujifilmImg.src}
                // comments={[
                //   { id: "1", author: "Alex", text: "I can pitch in $50!" },
                //   { id: "2", author: "Sam", text: "Great choice." }
                // ]}
              />
            </div>
          </div>

          {/* Right Column: Typography & Actions */}
          <div className="page__right-column">
            <h1 className="page__heading">
              Let your friends <br />
              know what <br />
              <span className="page__heading-italic">you're wishing for!</span>
            </h1>

            {/* Action Buttons & Inputs */}
            <div className="page__actions">
              
              {/* Create Wishlist Button */}
              <StylizedButton
                variant="solid"
                size="default"
                className="page__create-button"
                // onClick={handleCreateWishlist}
              >
                Create a Wishlist
              </StylizedButton>

              <div className="page__divider">
                <div className="page__divider-line"></div>
                <span className="page__divider-text">or</span>
                <div className="page__divider-line"></div>
              </div>

              {/* Guest Share Token Entry */}
              <form 
              // onSubmit={handleViewWishlist} 
              className="page__form">
                <label htmlFor="token" className="page__label">
                  View someone's wishlist
                </label>
                <div className="page__input-row">
                  <input
                    id="token"
                    type="text"
                    // value={shareToken}
                    // onChange={(e) => setShareToken(e.target.value)}
                    placeholder="Enter share token or URL..."
                    className="page__input"
                  />
                  <StylizedButton type="submit" className="page__view-button">
                    View
                  </StylizedButton>
                </div>
              </form>
              
            </div>
          </div>
        </div>
      </main>
  );
}

